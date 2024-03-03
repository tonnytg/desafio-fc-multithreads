# Desafio FC MultiThreads

Objetivo desse desafio é validar o entendimento sobre Multi Threads com Go


### Regras

Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
As duas requisições serão feitas simultaneamente para as seguintes APIs:
- - https://brasilapi.com.br/api/cep/v1/01153000 + cep
- - http://viacep.com.br/ws/" + cep + "/json/


**Os requisitos para este desafio são:**

- [X] Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- [X] O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- [X] Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.


### Como foi resolvido esse desafio?

Utilizei Go Routines para alimentar um Channel com os valores dentro de cep,
isso já começa dentro do contexto de Mutli Threads, pois cada Go Routine é uma Thread.

Depois criei uma função MakeRequest que vai iniciar duas Go Routines para fazer as requisições, porém há a regra que 
se uma das requisições retornar antes a outra deve ser cancelada.

Foi adicionado um select e um timeOut de 1 segundo caso não houve-se resposta no channel.

```
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, data)
```




### Output log

O teste do main é utilizando um slice de ceps coletados na internet

    ceps := []string{"01153000", "05541000", "08032330", "01550020", "03312006", "01415900", "05639050"}

Segue abaixo um exemplo de log esperado um pouco mais verboso

```
2024/03/02 21:26:46 Start Get CEP
2024/03/02 21:26:46 CEP 01153000
2024/03/02 21:26:46 Request 2 - URL https://brasilapi.com.br/api/cep/v1/01153000
2024/03/02 21:26:46 Request 1 - URL http://viacep.com.br/ws/01153000/json/
2024/03/02 21:26:47 Request 2 {"cep":"01153000","state":"SP","city":"São Paulo","neighborhood":"Barra Funda","street":"Rua Vitorino Carmilo","service":"correios-alt"}
2024/03/02 21:26:51 CEP 05541000
2024/03/02 21:26:51 Request 1 - URL http://viacep.com.br/ws/05541000/json/
2024/03/02 21:26:51 Request 2 - URL https://brasilapi.com.br/api/cep/v1/05541000
2024/03/02 21:26:52 Request 2 {"cep":"05541000","state":"SP","city":"São Paulo","neighborhood":"Jardim das Vertentes","street":"Avenida Albert Bartholome","service":"correios-alt"}
2024/03/02 21:26:56 CEP 08032330
2024/03/02 21:26:56 Request 1 - URL http://viacep.com.br/ws/08032330/json/
2024/03/02 21:26:56 Request 2 - URL https://brasilapi.com.br/api/cep/v1/08032330
2024/03/02 21:26:56 Request 2 {"cep":"08032330","state":"SP","city":"São Paulo","neighborhood":"Vila Nova Curuçá","street":"Rua Tachã","service":"correios-alt"}
2024/03/02 21:27:01 CEP 01550020
2024/03/02 21:27:01 Request 1 - URL http://viacep.com.br/ws/01550020/json/
2024/03/02 21:27:01 Request 2 - URL https://brasilapi.com.br/api/cep/v1/01550020
2024/03/02 21:27:02 Request 2 {"cep":"01550020","state":"SP","city":"São Paulo","neighborhood":"Vila Monumento","street":"Rua Doutor José Maria Azevedo","service":"correios-alt"}
```

Mesmo log abaixo, porém com a conexão WiFi desativada:
```
2024/03/02 21:27:06 CEP 03312006
2024/03/02 21:27:06 Request 2 - URL https://brasilapi.com.br/api/cep/v1/03312006
2024/03/02 21:27:06 Request 1 - URL http://viacep.com.br/ws/03312006/json/
2024/03/02 21:27:11 Timeout
2024/03/02 21:27:11 CEP 01415900
2024/03/02 21:27:11 Request 1 - URL http://viacep.com.br/ws/01415900/json/
2024/03/02 21:27:11 Request 2 - URL https://brasilapi.com.br/api/cep/v1/01415900
2024/03/02 21:27:16 Timeout
2024/03/02 21:27:16 CEP 05639050
2024/03/02 21:27:16 Request 1 - URL http://viacep.com.br/ws/05639050/json/
2024/03/02 21:27:16 Request 2 - URL https://brasilapi.com.br/api/cep/v1/05639050
2024/03/02 21:27:21 Timeout
2024/03/02 21:27:21 End Get CEP



```

