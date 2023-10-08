package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "net/http"

	_ "github.com/lib/pq"
)

type Usuario struct {
	id_usuario int
	username   string
	password   string
}

func getConnection() (*sql.DB, error) {
	connStr := "user=leolorenzo password=2190724 dbname=registros sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func crearUsuario(username, password string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	insertSQL := "INSERT INTO usuarios (username, password) VALUES ($1, $2);"
	_, err = db.Exec(insertSQL, username, password)
	if err != nil {
		return err
	}

	return nil
}

func obtenerUsuarioPorID(id_usuario int) (*Usuario, error) {
	db, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id_usuario, username, password FROM usuarios WHERE id_usuario = $1;"
	row := db.QueryRow(query, id_usuario)

	usuario := &Usuario{}
	err = row.Scan(&usuario.id_usuario, &usuario.username, &usuario.password)
	if err != nil {
		return nil, err
	}

	return usuario, nil
}

func actualizarUsuario(id_usuario int, username, password string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	updateSQL := "UPDATE usuarios SET username = $1, password = $2 WHERE id_usuario = $3;"
	_, err = db.Exec(updateSQL, username, password, id_usuario)
	if err != nil {
		return err
	}

	return nil
}

func eliminarUsuario(id_usuario int) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := "DELETE FROM usuarios WHERE id_usuario = $1;"
	_, err = db.Exec(deleteSQL, id_usuario)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := getConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Conexión a la base de datos establecida correctamente.")
	//WAOS
	// Definir un controlador para la página de formulario
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Mostrar el formulario HTML para crear un nuevo usuario
			// 	tmpl := `
			// <!DOCTYPE html>
			// <html>
			// <head>
			// 	<title>Formulario de Creación de Usuario</title>
			// </head>
			// <body>
			// 	<h1>Formulario de Creación de Usuario</h1>
			// 	<form method="post" action="/crear">
			// 		<label for="username">Nombre de usuario:</label>
			// 		<input type="text" name="username" required><br>
			// 		<label for="password">Contraseña:</label>
			// 		<input type="password" name="password" required><br>
			// 		<input type="submit" value="Crear Usuario">
			// 	</form>
			// </body>
			// </html>
			// `

			tmplHTML, err := template.ParseFiles("formulario.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmplHTML.Execute(w, nil)
		} else if r.Method == http.MethodPost {
			// Procesar el formulario cuando se envía
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Llamar a la función crearUsuario para insertar en la base de datos
			err := crearUsuario(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Redirigir o mostrar un mensaje de éxito
			http.Redirect(w, r, "/exito", http.StatusSeeOther)
		}
	})

	//FIN DEL WAOS
	http.ListenAndServe(":8080", nil)
}
