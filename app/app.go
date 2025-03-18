package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/vieitesss/ref/pages/references"
	"github.com/vieitesss/ref/pages/cheatsheet"
	"github.com/vieitesss/ref/pages/section"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	text string
}

func New() *Component {
	return &Component{
		mainRouter: router.New(),
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := section.New()

			return component, component.Init(section.Props{})
		},
		":sec/references": func(par router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := references.New()
			sec := par["sec"]

			return component, component.Init(references.Props{Section: sec})
		},
		":sec/references/:ref": func(par router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := cheatsheet.New()
			sec := par["sec"]
			ref := par["ref"]

			return component, component.Init(cheatsheet.Props{
				Section: sec,
				Reference: ref,
			})
		},
	})
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *Component) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}
