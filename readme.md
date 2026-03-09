## baixar tookit do IMB MQ
brew install ibm-messaging/ibmmq/mqdevtoolkit

## novas filas
QE.IMC.PB.MENSAGEM.FIX.ENTRADA.SESSAO30
QE.IMC.PB.MENSAGEM.FIX.ENTRADA.SESSAO70
QE.IMC.PB.MENSAGEM.FIX.SAIDA.SESSAO30
QE.IMC.PB.MENSAGEM.FIX.SAIDA.SESSAO70
QE.IMC.PB.MENSAGEM.FIX -- apagar
QE.IMC.PB.REPLICACAO.MENSAGENS.FIX

## exemplo de arquivo .env
```env
MQ_HOST=localhost(1414)
MQ_MANAGER=Q01
MQ_CHANNEL=CH.TO.Q01
```
obs.: executar com make

## executar
### Publish
> go run main.go publish QE.TEST abc

### Consume: 
> go run main.go consume QE.TEST