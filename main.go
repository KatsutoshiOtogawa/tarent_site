package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// api server
	//api_route := "http://127.0.0.1:8080/api/"

	// Renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer

	// Routes
	e.GET("/", hello)

	e.GET("/TarentDetail/:id", func(c echo.Context) error {

		id := c.Param("id")
		// 数字エラーチェック
		_, err := strconv.Atoi(id)
		if err != nil {
			return c.Render(http.StatusBadRequest, "404", nil)
		}

		db, err := sql.Open("sqlite3", os.Getenv("SQLITEPATH"))
		if err != nil {
			panic(err)
		}

		stmt, err := db.Prepare(`
		SELECT api_tarent.id,
		api_tarent.stage_name,
		api_tarent.family_name,
		api_tarent.first_name, 
		api_tarent.family_katakana_name,
		api_tarent.first_katakana_name,
		api_tarent.family_rome_name,
		api_tarent.first_rome_name,
		api_tarent.birth_date,     
		api_tarent.charm_point,
		api_tarent.image,
		replace(replace(api_tarent.review,'\r\n','<br />'),'\n','<br />') AS review,
		group_concat(DISTINCT api_tarentpersonality.id) AS tarent_personality_id,
		group_concat(DISTINCT api_tarentpersonality.name) AS tarent_personality_name,
		group_concat(DISTINCT api_tarentface.id) AS tarent_face_id,
		group_concat(DISTINCT api_tarentface.name) AS tarent_face_name,
		group_concat(DISTINCT api_tarentbody.id) AS tarent_body_id,
		group_concat(DISTINCT api_tarentbody.name) AS tarent_body_name,
		group_concat(DISTINCT api_tarentlowerbody.id) AS tarent_lowerbody_id,
		group_concat(DISTINCT api_tarentlowerbody.name) AS tarent_lowerbody_name,
		group_concat(DISTINCT api_tarentupperbody.id) AS tarent_upperbody_id,
		group_concat(DISTINCT api_tarentupperbody.name) AS tarent_upperbody_name,
		api_tarentbrasize.id AS tarent_brasize_id,
		api_tarentbrasize.name AS tarent_brasize_name,
		group_concat(DISTINCT api_tarentart.id) AS tarent_art_id,
		group_concat(DISTINCT api_tarentart.name) AS tarent_art_name,
		group_concat(DISTINCT api_tarentsite.url) AS tarent_site_url,
		group_concat(DISTINCT api_sitetype.name) AS site_type_name,
		group_concat(DISTINCT CASE api_sitetype2.name
			WHEN 'twitter' then api_tarentinfositeembed.html
			ELSE NULL
		END) AS  twitter_embed_html,
		group_concat(DISTINCT CASE api_sitetype2.name
			WHEN 'instagram' then api_tarentinfositeembed.html
			ELSE NULL
		END) AS  instagram_embed_html,
		group_concat(DISTINCT CASE api_sitetype2.name
			WHEN 'youtube' then api_tarentinfositeembed.html
			ELSE NULL
		END) AS  youtube_embed_html
	
		FROM api_Tarent
		INNER JOIN api_tarent_tarent_personality
			ON api_Tarent.id = api_tarent_tarent_personality.tarent_id
		INNER JOIN api_tarentpersonality
			ON api_tarent_tarent_personality.tarentpersonality_id = api_tarentpersonality.id
		
		INNER JOIN api_tarent_tarent_face
			ON api_Tarent.id = api_tarent_tarent_face.tarent_id
		INNER JOIN api_tarentface
			ON api_tarent_tarent_face.tarentface_id = api_tarentface.id
			
		INNER JOIN api_tarent_tarent_body
			ON api_Tarent.id = api_tarent_tarent_body.tarent_id
		INNER JOIN api_tarentbody
			ON api_tarent_tarent_body.tarentbody_id = api_tarentbody.id
		
		INNER join api_tarent_tarent_lower_body
			ON api_Tarent.id= api_tarent_tarent_lower_body.tarent_id
		INNER JOIN api_tarentlowerbody
			ON api_tarent_tarent_lower_body.tarentlowerbody_id = api_tarentlowerbody.id
			
		INNER join api_tarent_tarent_upper_body
			ON api_Tarent.id= api_tarent_tarent_upper_body.tarent_id
		INNER JOIN api_tarentupperbody
			ON api_tarent_tarent_upper_body.tarentupperbody_id = api_tarentupperbody.id
		INNER JOIN api_tarentbrasize
			ON api_Tarent.tarent_bra_size_id = api_tarentbrasize.id
		INNER JOIN api_tarentart
			ON api_tarent.id = api_tarentart.tarent_id
		INNER JOIN api_tarentsite
			ON api_tarent.id = api_tarentsite.tarent_id
		INNER JOIN api_sitetype
			ON api_tarentsite.site_type_id = api_sitetype.id
	
		INNER JOIN api_tarentinfositeembed
			ON api_tarent.id = api_tarentinfositeembed.tarent_id
		INNER JOIN api_sitetype AS api_sitetype2
			ON api_tarentinfositeembed.site_type_id = api_sitetype2.id
			AND api_sitetype2.name IN ('instagram','youtube','twitter')
	
		WHERE api_Tarent.id = ?
		GROUP BY api_tarent.stage_name,
			api_tarent.family_name,
			api_tarent.first_name, 
			api_tarent.family_katakana_name,
			api_tarent.first_katakana_name,
			api_tarent.family_rome_name,
			api_tarent.first_rome_name,
			api_tarent.image,
			api_tarent.review,
			api_tarent.birth_date,     
			api_tarent.charm_point,
			api_tarentbrasize.id,
			api_tarentbrasize.name`)
		if err != nil {
			panic(err)
		}
		rows, err := stmt.Query(id)
		if err != nil {
			panic(err)
		}
		type Data struct {
			Content_id                      int
			Content_stage_name              string
			Content_family_name             string
			Content_first_name              string
			Content_family_katakana_name    string
			Content_first_katakana_name     string
			Content_family_rome_name        string
			Content_first_rome_name         string
			Content_birth_date              string
			Content_charm_point             string
			Content_image                   string
			Content_review                  string
			Content_tarent_personality_id   string
			Content_tarent_personality_name string
			Content_tarent_face_id          string
			Content_tarent_face_name        string
			Content_tarent_body_id          string
			Content_tarent_body_name        string
			Content_tarent_upperbody_id     string
			Content_tarent_upperbody_name   string
			Content_tarent_lowerbody_id     string
			Content_tarent_lowerbody_name   string
			Content_tarent_brasize_id       int
			Content_tarent_brasize_name     string
			Content_tarent_art_id           string
			Content_tarent_art_name         string
			Content_tarent_site_url         string
			Content_site_type_name          string
			Content_twitter_embed_html      string
			Content_instagram_embed_html    string
			Content_youtube_embed_html      string
		}
		var data Data
		for rows.Next() {
			err = rows.Scan(
				&(data.Content_id),
				&(data.Content_stage_name),
				&(data.Content_family_name),
				&(data.Content_first_name),
				&(data.Content_family_katakana_name),
				&(data.Content_first_katakana_name),
				&(data.Content_family_rome_name),
				&(data.Content_first_rome_name),
				&(data.Content_birth_date),
				&(data.Content_charm_point),
				&(data.Content_image),
				&(data.Content_review),
				&(data.Content_tarent_personality_id),
				&(data.Content_tarent_personality_name),
				&(data.Content_tarent_face_id),
				&(data.Content_tarent_face_name),
				&(data.Content_tarent_body_id),
				&(data.Content_tarent_body_name),
				&(data.Content_tarent_upperbody_id),
				&(data.Content_tarent_upperbody_name),
				&(data.Content_tarent_lowerbody_id),
				&(data.Content_tarent_lowerbody_name),
				&(data.Content_tarent_brasize_id),
				&(data.Content_tarent_brasize_name),
				&(data.Content_tarent_art_id),
				&(data.Content_tarent_art_name),
				&(data.Content_tarent_site_url),
				&(data.Content_site_type_name),
				&(data.Content_twitter_embed_html),
				&(data.Content_instagram_embed_html),
				&(data.Content_youtube_embed_html),
			)
			if err != nil {
				panic(err)
			}
		}

		return c.Render(http.StatusOK, "TarentDetail", data)
	})

	// example
	e.GET("/page1", func(c echo.Context) error {
		// テンプレートに渡す値

		data := struct {
			Content_a string
			Content_b string
			Content_c string
			Content_d string
		}{
			Content_a: "雨が降っています。",
			Content_b: "明日も雨でしょうか。",
			Content_c: "台風が近づいています。",
			Content_d: "Jun/11/2018",
		}
		return c.Render(http.StatusOK, "page1", data)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
