# Banco de dados e bibliotecas utilizadas
 - Banco de dados utilizado: Postgres
 - Biblioteca para router: gorilla/mux
 - Biblioteca para migration: golang-migrate

# Estrutra de pastas
 - config: module responsável por carregar as variáveis de ambiente. 
 - database: module (db_service.go) responsável por conectar-se com o banco de dados Postgres, expor uma conexão global ("Conn *sql.DB") e rodar o migrations para criação dos objetos.
 - handlers: module com os controllers e adição dos mesmos ao router.
 - interfaces: especificação das interfaces a serem implementadas.
 - migration: arquivos de sql para rodar o migration do banco de dados.
 - mocks: estruturas auxiliares para mockar repositórios, serviços e http.Client nos testes unitários.
 - models: estruturas que representam os dados dos jsons nos requests e responses.
 - services: camada de serviço acionada nos controllers para executar os procedimento necessários.
 - repositories: camada de acesso direto ao banco de dados (operações DML)

 OBS: Todas as estruturas de serviços (module services) e repositórios (module repositories) implementam as interfaces especificadas no modules interfaces, facilitando assim a injeção de dependências e criação dos testes unitários.

# Para execução de testes unitários
  rodar: go test ./services/...
  arquivo com as impletações dos testes: services/freigth_test.go


# Para execução do projeto via docker

* Necessário configurar as variáveis de ambiente no arquivo docker-compose.yml conforme os valores passados
no arquivo html de instrução desse desafio:

   - REG_NUMBER= "CNPJ Remetente - somente números"
   - TOKEN= "Tojen de autenticação" 
   - SYS_CODE= "Código da plataforma"
   - DISP_ZIPCODE= "Cep - somente números"
  
* Acionar o docker-compose para subir a imagem da API e do banco de dados:
   docker-compose up -d  

# Exemplo de JSON retornado no Endpoint GET(Quotes)

```json
{
	"carrier": [
		{
			"name": "CORREIOS", //nome da transportadora
			"total_freight": 642.03, //total frete
			"average_freight": 321.015, //média frete
			"qty_results": 2 // quantidade de resultados
		},
		{
			"name": "CORREIOS - SEDEX",
			"total_freight": 642.03,
			"average_freight": 321.015,
			"qty_results": 2
		}
	],
	"min_price": 264.48, //frete mais barato
	"max_price": 377.55 //frete mais caro
}










