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

type DatosPagina struct {
	Usuarios []Usuario
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

func actualizarUsuario(id_usuario string, username, password string) error {
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

func eliminarUsuario(id_usuario string) error {
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

func obtenerTodosLosUsuarios() ([]Usuario, error) {
	db, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id_usuario, username FROM usuarios;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var usuario Usuario
		err := rows.Scan(&usuario.id_usuario, &usuario.username)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

var plantilla = template.Must(template.ParseGlob("formulario.html"))

func main() {
	db, err := getConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Conexión a la base de datos establecida correctamente.")
	usuarios, err := obtenerTodosLosUsuarios()
	if err != nil {
		log.Fatal(err)
	}

	// Mostrar los usuarios en la consola
	for _, usuario := range usuarios {
		fmt.Printf("ID: %d, Usuario: %s\n", usuario.id_usuario, usuario.username)
	}

	//WAOS
	// Definir un controlador para la página de formulario
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

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
			fmt.Println("\nUsuario Registrado!!\n")
			// Redirigir o mostrar un mensaje de éxito
			http.Redirect(w, r, "/", http.StatusSeeOther)

		}

	})

	http.HandleFunc("/actualizar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		id_usuario := r.FormValue("id_usuario")
		username := r.FormValue("username")
		password := r.FormValue("n_password")
		// Convertir la cadena del ID en un entero

		// Llamar a la función crearUsuario para insertar en la base de datos
		err := actualizarUsuario(id_usuario, username, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("\nUsuario Actualizado!!\n")
		// Redirigir o mostrar un mensaje de éxito
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/eliminar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		id_usuario := r.FormValue("id_usuario")

		err := eliminarUsuario(id_usuario)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("\nUsuario Eliminado!!\n")
		// Redirigir o mostrar un mensaje de éxito
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	//FIN DEL WAOS
	http.ListenAndServe(":8080", nil)
}
