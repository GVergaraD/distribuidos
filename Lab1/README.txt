Jean Aravena 201673573-9
Gabriel Vergara 201673605-0 

Los archivos fueron clonados desde un repositorio(github) con git clone
Para ejecturar los archivos es necesario hacer lo siguiente:
	-dist81->Logistica
	desde carpeta distribuidos/Lab/Src/ 
	correr ./server

	-dist82->Finanzas
	desde carpeta distribuidos/Lab/Src/
	correr ./financiero

	-dist83->Cliente
	desde carpeta distribuidos/Lab/Src/
	correr ./client

	-dist84->Camiones
	desde carpeta distribuidos/Lab/Src/
	correr ./camiones

Consideraciones:
-Logistica y finanzas deben ser los primeros en ser ejecutados
-Los numero de seguimiento asociado a las pymes va parte desde N1
-Los id de los paquetes parte de 1
-Para consultar el seguimiento abrir una nueva consola si se esta enviando paquetes desde el cliente
-Se asume que la persona que va a usar el sistema va a ingresar datos validos con por consola
-El enunciado de la tarea dice que retail tiene 3 intentos maximo y pymes tiene 2 reintentos, por lo cual se asume que ambos tienen como maximo 3 intentos
-En finanzas se asume que perdidas se asocia a cada reintento y ganancia se asocia a lo que se puede ganar con cada paquete
-Si se quiere ejecutar nuevamente el sistema borrar todo lo que encuentra en cada csv (correrlo de cero)


En caso de error por GOPATH o GOROOT ingresar comando:export PATH="$PATH:$(go env GOPATH)/bin"
en el directorio distribuidos/Lab/Src 