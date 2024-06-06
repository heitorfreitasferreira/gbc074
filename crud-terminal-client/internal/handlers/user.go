package handlers

/*
	message Usuario{
	// CPF do usuario (chave)
	string cpf     = 1;
	string nome    = 2;
	// campo presente apenas no portal biblioteca
	bool bloqueado = 3;
}
*/

type user struct {
	cpf       string
	nome      string
	bloqueado bool
}
