package main

import (
    "github.com/gotk3/gotk3/gtk"
    "log"
)

type Menu struct {
    entries []MenuEntry
}

func (m *Menu) AddEntryWithAction(label string, next *Menu, action func()) {
    if next == nil { next = m }

    m.entries = append(m.entries, MenuEntry{
        Label: label,
        Next: next,
        Action: action,
    })
}

func (m *Menu) AddEntry(label string, next *Menu) {
    m.AddEntryWithAction(label, next, nil)
}

type MenuEntry struct {
    Label string
    Next *Menu
    Action func()
}

func (e MenuEntry) Use() *Menu {
    if e.Action != nil { e.Action() }

    return e.Next
}

func (m Menu) ProcessNextMenu(box *gtk.Box) {
    box.GetChildren().Foreach(func (child any) {
        btn, _ := child.(*gtk.Widget)
        btn.Destroy()
    })

    for _, entry := range m.entries {
        btn, err := gtk.ButtonNewWithLabel(entry.Label)

        if err != nil {
            log.Panic(err)
        }

        currentEntry := entry
        btn.Connect("clicked", func() {
            currentEntry.Use().ProcessNextMenu(box)
        })

        box.Add(btn)
    }

    box.ShowAll()
}

func (m Menu) GtkWidget() *gtk.Widget {
    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)

    if err != nil {
        log.Panic(err)
    }

    m.ProcessNextMenu(box)

    return &box.Widget
}

func makeMainMenu(info *gtk.Label) *Menu {
    var mainMenu, newGameMenu, optionsMenu Menu

    var gameResult bool

    playGame := func() {
        if gameResult {
            info.SetText("Вы выиграли!")
        } else {
            info.SetText("Вы проиграли!")
        }
    }

    clearInfo := func() { info.SetText("") }

    mainMenu.AddEntryWithAction("Новая игра", &newGameMenu, playGame)

    mainMenu.AddEntry("Настройки", &optionsMenu)
    mainMenu.AddEntryWithAction("Выйти", nil, gtk.MainQuit)

    newGameMenu.AddEntryWithAction("Начать заново", nil, playGame)
    newGameMenu.AddEntryWithAction("Выйти в главное меню", &mainMenu, clearInfo)

    optionsMenu.AddEntryWithAction("Хочу всегда выигрывать", &mainMenu, func() {
        gameResult = true
    })

    optionsMenu.AddEntryWithAction("Хочу всегда проигрывать", &mainMenu, func() {
        gameResult = false
    })

    return &mainMenu
}

func main() {
    gtk.Init(nil)

    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

    if err != nil {
        log.Fatal(err)
    }

    win.Connect("destroy", func() {
        gtk.MainQuit()
    })

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)

    if err != nil {
        log.Panic(err)
    }

    box.SetMarginTop(24)
    box.SetMarginBottom(24)
    box.SetMarginStart(24)
    box.SetMarginEnd(24)

    win.Add(box)

    label, err := gtk.LabelNew("")

    if err != nil {
        log.Panic(err)
    }

    box.Add(label)
    box.Add(makeMainMenu(label).GtkWidget())

    win.ShowAll()

    gtk.Main()
}
