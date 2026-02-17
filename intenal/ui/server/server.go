package server

import (
	"andrew_chat/intenal/debug"
	"andrew_chat/intenal/domain"
	"andrew_chat/intenal/server"
	"andrew_chat/intenal/ui"
	"andrew_chat/intenal/ui/keys"
	"andrew_chat/intenal/ui/types"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// =============================================================================
// Server
// =============================================================================

// implements bubbletea.model
type ServerModel struct {
	list  	*ui.List
	width    int
	height   int
	ss       *server.ServerService
}

func NewServer(servers []domain.Server) *ServerModel {
	return &ServerModel{
		ss:      server.NewServerService(),
	}
}

func serversToItems(servers []domain.Server) []list.Item {
    items := make([]list.Item, len(servers))
    for i := range servers {
        items[i] = servers[i] 
    }
    return items
}


var fields = []types.InputFieldSpec{
    {
        Name:        "name",
        Title:       "Server Name",
        Desc:        "Display name for the server",
        Placeholder: "Production Server EU",
    },
    {
        Name:        "address",
        Title:       "Address",
        Desc:        "Server IP address or hostname",
        Placeholder: "192.168.1.10 or example.com",
    },
    {
        Name:        "port",
        Title:       "Port",
        Desc:        "Connection port",
        Placeholder: "4567",
    },
    {
        Name:        "username",
        Title:       "Username",
        Desc:        "User for authentication",
        Placeholder: "admin",
    },
    {
        Name:        "protocol",
        Title:       "Protocol",
        Desc:        "Connection protocol",
        Placeholder: "ssh, http, https",
    },
}

func (m *ServerModel) makeInputForm() tea.Model {
    prompt := ">> "
    form := ui.NewInputFormModel(prompt, fields)
    return form
}


func (m *ServerModel) initOptions(serverID string) []ui.Option{
	srvItem := m.list.SelectedItem()
	if srvItem == nil {
		panic("nil selected item")
	}

	srv := srvItem.(domain.Server)

	return []ui.Option{
		{
			Name: "delete",
			Action: func() tea.Cmd {
				var cmds []tea.Cmd
				err := m.ss.Remove(serverID)
				cmds = append(cmds, ui.NewDeleteCmd(m))
				if err != nil {
					cmds = append(cmds, ui.NewErrCmd("remove failed"))
				}
				return tea.Batch(cmds...)
			},
		},
		{
			Name: "connect",
			Action: func() tea.Cmd {
				var cmds []tea.Cmd
				err := m.ss.Connect(srv)
				if err != nil {
					cmds = append(cmds, ui.NewErrCmd("connect failed"))
				}
				return tea.Batch(cmds...)
			},
		},
		{
			Name: "update",
			Action: func() tea.Cmd {
				var cmds []tea.Cmd
				srv.Name = srv.Name + " (updated)"
				err := m.ss.Update(srv)
				if err != nil {
					cmds = append(cmds, ui.NewErrCmd("update failed"))
				}
				return tea.Batch(cmds...)
			},
		},
	}
}

func (m *ServerModel) Init() tea.Cmd {
	items := serversToItems(m.ss.GetServers())
	m.list = ui.NewList(items, list.NewDefaultDelegate(), m.width, m.height)
	return nil
}

func (m *ServerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.list.Update(msg)
	

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Choose):
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				return m, cmd
			}
			srv := selectedItem.(domain.Server)
			srvMsg := tea.Cmd(func() tea.Msg {
				err := m.ss.Connect(srv)
				return types.ServerMsg{
					Status:  server.StatusConnected,
					Server:  srv,
					Success: err == nil,
				}
			})

			inputModel := m.makeInputForm()
			navMsg := ui.NewCreateCmd(types.PositionTopRight, inputModel)
			return m, tea.Batch(srvMsg, navMsg)
		case key.Matches(msg, keys.Keys.Close):
			panic("unexpected sequence")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	
	case types.InputFormMsg:
		server, err := formMsgToServer(msg)
		if err != nil {
			panic("wrong format")
		}

		if err = m.ss.Add(server); err != nil{
			panic(err)
		}
		items := serversToItems(m.ss.GetServers())
		m.list = ui.NewList(items, list.NewDefaultDelegate(), m.width, m.height)
		debug.DebugDump(debug.V, "types.InputFormMsg", server)
	}

	return m, cmd
}

// View
func (m *ServerModel) View() string {
    // panic(fmt.Sprintf("START/%s/END", m.listview.View()))
	return m.list.View()
}

func formMsgToServer(msg types.InputFormMsg) (domain.Server, error) {
    var server domain.Server

    for _, field := range msg.Fields {
        switch field.Name {
        case "name":
            server.Name = field.Value
        case "address":
            server.Address = field.Value
        case "port":
            p, err := strconv.Atoi(field.Value)
            if err != nil {
                return server, fmt.Errorf("invalid port: %v", err)
            }
            server.Port = p
        case "protocol":
            server.Protocol = field.Value
        case "username":
            server.Username = field.Value
        case "description":
            server.Desc = field.Value
        case "environment":
            server.Environment = field.Value
        case "region":
            server.Region = field.Value
        case "active":
            if field.Value == "true" || field.Value == "1" {
                server.Active = true
            } else {
                server.Active = false
            }
        }
    }

    return server, nil
}
