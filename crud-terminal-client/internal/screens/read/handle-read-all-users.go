package screens

type ReadAllUsers struct {
}

// func (r *ReadAllUsers) HandleReadAllUsers() tea.Msg {
// 	conn := utils.GetConn(50051)
// 	defer conn.Close()

// 	client := br_ufu_facom_gbc074_projeto_cadastro.NewPortalCadastroClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	response, err := client.ObtemTodosUsuarios(ctx, &br_ufu_facom_gbc074_projeto_cadastro.Vazia{})
// 	if err != nil {
// 		log.Fatalf("Não foi possível ler os usuários: %v", err)
// 	}
// 	// var us []br_ufu_facom_gbc074_projeto_cadastro.Usuario = response.get

// 	for _, user := range users {
// 		log.Printf("CPF: %s, Nome: %s\n", user.Cpf, user.Nome)
// 	}

// }
