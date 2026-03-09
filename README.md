## baixar tookit do IMB MQ
brew install ibm-messaging/ibmmq/mqdevtoolkit

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