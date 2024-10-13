# Rate Limiting - Go Graduation Project 

## Visão Geral

O  (ou limitador de taxa) é uma técnica essencial para controlar e restringir o número de requisições que um cliente—seja um usuário, IP ou serviço—pode realizar em uma aplicação em um intervalo de tempo específico.

Neste projeto, implementamos um *rate limiting* em uma aplicação simples de livraria, com o objetivo de *demonstrar a eficácia* dessa técnica e sua *relevância em sistemas reais*, especialmente na *mitigação de ataques DDoS* e no controle do uso abusivo de recursos. O projeto foi desenvolvido *totalmente em Go*, utilizando *Docker* e *Redis* para gerenciamento de cache. A aplicação se conecta a um banco de dados *MySQL*, e toda a arquitetura do sistema foi projetada com base nos princípios da *Clean Architecture*, o que assegura *escalabilidade* e facilidade de manutenção.

Para a implementação do *servidor e roteador HTTP*, utilizamos a biblioteca *Chi*, enquanto o gerenciamento de injeções de dependências foi feito com o *Google Wire*, garantindo uma estrutura de código limpa e modular.

Esse projeto é um trabalho de conclusão de curso da pós graduação Go Expert.

## Uso

Para utilizar o projeto *Rater Limit*, siga os passos abaixo:

### Pré-requisitos

Certifique-se de ter o Docker e o Docker Compose instalados em seu sistema.

### Clonando o Repositório

Clone o repositório do GitHub:

```bash
git clone https://github.com/guirialli/rater-limit.git
cd rater-limit
```

### Configuração do Ambiente

Antes de iniciar a aplicação, você pode configurar as variáveis de ambiente no arquivo `.env`, localizado em `cmd/rater_limit`. As configurações incluem limites de taxa, informações do banco de dados, segurança JWT e configurações do Redis. Aqui estão algumas das variáveis importantes:

- *Rate Limiting*:

  - `IP_REFRESH_ACCESS`: Tempo (em segundos) para resetar as tentativas de acesso por IP.
  - `IP_TRYS_MAX`: Limite de tentativas por IP antes do bloqueio.
  - `JWT_REFRESH_ACCESS`: Tempo (em segundos) para resetar as tentativas de acesso por JWT.
  - `JWT_TRYS_MAX`: Limite de tentativas por JWT antes do bloqueio.
  - `BLOCK_TIMEOUT`: Tempo (em minutos) que um IP ou JWT permanecerá bloqueado após exceder o limite de tentativas.

- *Banco de Dados*:

  - `DB_HOSTNAME`: Nome do host ou container que executa o MySQL.
  - `DB_PORT`: Porta do MySQL.
  - `DB_DATABASE`: Nome do banco de dados.
  - `DB_USER`: Usuário do banco de dados.
  - `DB_PASSWORD`: Senha do banco de dados.

- *Configurações do Redis*:

  - `REDIS_ADDR`: Endereço e porta do servidor Redis.

  - `REDI``S_PASSWORD`: Senha para autenticação no Redis.

```
    # Rater Limiting
    
    # IP_REFRESH_ACCESS: Tempo em *segundos (s)* para resetar as tentativas de acesso de IPs ao limite máximo.
    
    IP_REFRESH_ACCESS=1
    
    # IP_TRYS_MAX: Número máximo de tentativas permitidas por IP antes do bloqueio.
    
    IP_TRYS_MAX=5
    
    # JWT_REFRESH_ACCESS: Tempo em *segundos (s)* para resetar as tentativas de acesso de JWTs ao limite máximo.
    
    JWT_REFRESH_ACCESS=1
    
    # JWT_TRYS_MAX: Número máximo de tentativas permitidas por JWT antes do bloqueio.
    
    JWT_TRYS_MAX=10
    
    # BLOCK_TIMEOUT: Tempo em *minutos (m)* que um IP ou JWT permanecerá bloqueado após exceder o limite de tentativas.
    
    BLOCK_TIMEOUT=1
    
    # Database
    
    # DB_HOSTNAME: Nome do host ou container que executa o banco de dados MySQL.
    
    DB_HOSTNAME="rater_limit_db"
    
    # DB_PORT: Porta na qual o banco de dados MySQL está escutando.
    
    DB_PORT="3306"
    
    # DB_DATABASE: Nome do banco de dados utilizado para armazenar as informações de rate limiting.
    
    DB_DATABASE="rater_limit"
    
    # DB_USER: Usuário do banco de dados.
    
    DB_USER="root"
    
    # DB_PASSWORD: Senha do usuário do banco de dados.
    
    DB_PASSWORD="root"
    
    # Security JWT
    
    # JWT_EXPIRE_IN: Duração do token JWT antes da expiração (unidade definida em JWT_UNIT_TIME).
    
    JWT_EXPIRE_IN="100"
    
    # JWT_SECRET: Chave secreta usada para assinar e verificar os tokens JWT.
    
    JWT_SECRET="test"
    
    # JWT_UNIT_TIME: Unidade de tempo usada na expiração do JWT:
    
    # 's' = segundos, 'm' = minutos, 'h' = horas, 'd' = dias.
    
    JWT_UNIT_TIME="h"
    
    # Webserver
    
    # IP: Endereço IP que o servidor irá escutar (deixe em branco para escutar em todos os endereços).
    
    IP=""
    
    # PORT: Porta em que o servidor web estará disponível.
    
    PORT="8080"
    
    # Redis
    
    # REDIS_ADDR: Endereço e porta do servidor Redis.
    
    REDIS_ADDR="redis:6379"
    
    # REDIS_PASSWOR: Senha para autenticação no Redis (deixe em branco se não houver).
    
    REDIS_PASSWOR=""
    
    # REDIS_DB: Número do banco Redis a ser utilizado.
    
    REDIS_DB=0
```

    

### Subindo a Aplicação

Para iniciar a aplicação, utilize o comando:

```bash
docker-compose up --build
```

### Acessando a API

Após a aplicação estar em execução, você pode acessar a API em `http://localhost:8080`.

### Swagger

ara facilitar a utilização e a documentação da API, foi configurado o Swagger. É importante destacar que, por padrão, o *Rate Limiting* não se aplica ao Swagger, permitindo que desenvolvedores e usuários explorem a documentação da API sem restrições. O Swagger pode ser acessado através da rota `http://127.0.0.1:8080/swagger/`. O Rater Limit irá ignorar qualquer acesso ao swagger permitindo assim a livre consulta da api.

Além disso, as rotas sob o prefixo `/public/*` não requerem autenticação, permitindo que usuários realizem testes de acesso tanto por IP quanto por JWT. Isso facilita a verificação do funcionamento do sistema de *Rate Limiting* sem a necessidade de estar logado, promovendo um ambiente mais acessível para testes e desenvolvimento.

### Postman

Para facilitar os testes da aplicação, foi incluído um arquivo de coleção do Postman localizado em `test/postman/Rater Limit.postman_collection.json`. Este arquivo contém várias requisições pré-configuradas que permitem simular diferentes cenários de uso do rate limiting, tanto para IP quanto para JWT.

Com essa coleção, os desenvolvedores podem facilmente testar a funcionalidade de rate limiting da aplicação, verificando como ela se comporta sob diferentes condições e garantindo que as regras de bloqueio e limite de tentativas estejam funcionando corretamente. Para usar a coleção, basta importar o arquivo no Postman e executar as requisições conforme necessário, facilitando o processo de validação e testes da aplicação.
