package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "library.db")
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func createTable() error {
	_, err := db.Exec(`
 		CREATE TABLE IF NOT EXISTS books (
 			id INTEGER PRIMARY KEY AUTOINCREMENT,
 			book_barcode_number TEXT
 			days TEXT
		);
	`)
	return err
}

func insertBook(bookNumber string) error {
	_, err := db.Exec("INSERT INTO books (book_barcode_number) VALUES (?)", bookNumber)
	return err
}

func makeUI(tabs *container.AppTabs) *fyne.Container {
	return container.New(layout.NewGridLayout(1),
		tabs)
}

func makeUI2(entry *widget.Entry, button *widget.Button) *fyne.Container {
	return container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.NewVBox(entry),
		container.NewVBox(button),
	)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Kütüphane")
	myWindow.Resize(fyne.NewSize(300, 400))
	// if err := os.Chdir(myApp.Storage().RootURI().Path()); err != nil {
	// 	log.Fatal(err)
	// }

	var err error
	db, err = openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable()
	if err != nil {
		log.Fatal(err)
	}

	input := widget.NewEntry()
	input.SetPlaceHolder("Kitap Numarasını Girin")

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Ana Sayfa", theme.HomeIcon(), container.NewVBox(widget.NewButton("Kitap Kaydet", func() {
			w2 := myApp.NewWindow("Kütüphane")
			w2.Resize(fyne.NewSize(300, 400))
			w2.SetContent(makeUI2(input, widget.NewButton("Kaydet", func() {
				bookNumber := input.Text
				if bookNumber == "" {
					fyne.LogError("You Cannot Enter Null Barcode Number", err)
				} else {
					err := insertBook(bookNumber)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Book with number %s inserted successfully\n", bookNumber)
				}
			})))
			w2.Show()
			log.Print("Clicked")
		}))),
		container.NewTabItem("Ayarlar", container.NewVBox(widget.NewSelect([]string{"Türkçe", "Arapça"}, func(s string) {

		}))),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	myWindow.SetContent(makeUI(tabs))

	myWindow.Show()
	myApp.Run()
}
