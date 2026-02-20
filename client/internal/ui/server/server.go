package server

import (
	"andrew_chat/client/internal/domain"
	"andrew_chat/client/internal/server"
	"andrew_chat/client/internal/ui"
	"andrew_chat/client/internal/ui/keys"
	"andrew_chat/client/internal/ui/types"
	"encoding/json"
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
	//main list with all servers
	list *ui.List

	//another window that is opened and controlled by ServerModel
	descModel *ui.TextView

	width  int
	height int
	ss     *server.ServerService
}

func NewServer(servers []domain.Server) *ServerModel {
	return &ServerModel{
		ss: server.NewServerService(),
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
	form := ui.NewInputFormModel(prompt, fields, m.inputFormAction)
	return form
}

func (m *ServerModel) newServerMsg(srv domain.Server, status int) tea.Msg {
	return types.ServerMsg{
		Status: status,
		Server: srv,
	}
}

func (m *ServerModel) initOptions(serverID string) []ui.Option {
	srvItem := m.list.SelectedItem()
	if srvItem == nil {
		panic("nil selected item")
	}

	srv := srvItem.(domain.Server)

	return []ui.Option{
		{
			Name: "connect",
			Action: func() tea.Cmd {
				return tea.Sequence(
					func() tea.Msg {
						return m.newServerMsg(srv, server.StatusConnecting)
					},

					func() tea.Msg {
						err := m.ss.Connect(srv)
						if err != nil {
							return m.newServerMsg(srv, server.StatusDisconnected)
						}
						return m.newServerMsg(srv, server.StatusConnected)
					},

					func() tea.Msg {
						m.updateList()
						return nil
					},
				)
			},
		},
		{
			Name: "delete",
			Action: func() tea.Cmd {
				var cmds []tea.Cmd
				err := m.ss.Remove(serverID)
				if err != nil {
					cmds = append(cmds, ui.NewErrCmd("remove failed"))
				}
				m.updateList()
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

func (m *ServerModel) marshalSelectedItem() []byte {
	b, err := json.MarshalIndent(m.list.SelectedItem(), "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

func (m *ServerModel) createServerTextWinCmd() tea.Cmd {
	b := m.marshalSelectedItem()
	listview := ui.NewTextView()
	listview.SetContent(string(b))
	m.descModel = listview
	return ui.NewCreateCmd(types.PositionTopRight, listview, false)
}

func (m *ServerModel) Init() tea.Cmd {
	m.updateList()

	return m.createServerTextWinCmd()
}

func (m *ServerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// update list
	_, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Choose):
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				return m, tea.Batch(cmds...)
			}
			srv := selectedItem.(domain.Server)
			opts := m.initOptions(srv.ID)
			control := ui.NewControlPane(opts)
			navCmd := ui.NewCreateCmd(types.PositionBotLeft, control, true)
			cmds = append(cmds, navCmd)
			return m, tea.Batch(cmds...)
		case key.Matches(msg, keys.Keys.Down) || key.Matches(msg, keys.Keys.Up):
			b := m.marshalSelectedItem()
			m.descModel.SetContent(string(b))
			// _, cmd = m.descModel.Update(msg)
			cmds = append(cmds, cmd)
			// case key.Matches(msg, keys.Keys.Close):
			// 	panic("unexpected sequence")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

// View
func (m *ServerModel) View() string {
	// panic(fmt.Sprintf("START/%s/END", m.listview.View()))
	return m.list.View()
}

func (m *ServerModel) updateList() {
	items := serversToItems(m.ss.GetServers())
	m.list = ui.NewList(items, list.NewDefaultDelegate(), m.width, m.height)
}

func (m *ServerModel) inputFormAction(values []types.InputFieldValue) tea.Cmd {
	var server domain.Server

	for _, field := range values {
		switch field.Name {
		case "name":
			server.Name = field.Value
		case "address":
			server.Address = field.Value
		case "port":
			p, err := strconv.Atoi(field.Value)
			if err != nil {
				panic("invalid port")
				// return server, fmt.Errorf("invalid port: %v", err)
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
	if err := m.ss.Add(server); err != nil {
		panic(err.Error())
	}

	m.updateList()
	return nil
}
