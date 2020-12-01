Jean Aravena 201673573-9
Gabriel Vergara 201673605-0 

Los archivos fueron clonados desde un repositorio(github) con git clone

El programa se encuentra distribuido en 4 maquinas en donde la primera cuenta con 1 datanode y 1 cliente, la 2da y 3ra cuentan con los datanodes 2 y 3 respectivamente y la 4ta cuenta con un namenode que lleva el Log de los libros distribuidos.
Para ejecturar los archivos es necesario hacer lo siguiente:
	-dist81->Datanode1
	desde carpeta distribuidos/Lab2/Src/
	correr ./datanode1.go

	-dist81->Cliente
	desde carpeta distribuidos/Lab2/Src/Client 
	correr ./client.go

	-dist82->Datanode2
	desde carpeta distribuidos/Lab2/Src/DataNode2
	correr ./datanode2.go

	-dist83->Datanode3
	desde carpeta distribuidos/Lab2/Src/DataNode3
	correr ./datanode3.go

	-dist84->Namenode
	desde carpeta distribuidos/Lab2/Src/NameNode
	correr ./namenode.go


Consideraciones:
-Los libros utilizados estan en formato pdf.(https://www.elejandria.com/coleccion/libros-llevados-al-cine)
-El Log se escribe en un archivo csv
-Los libros disponibles se encuentran en la carpeta Libros
-En ninguno de los 2 métodos se implementó el algoritmo para el caso en que acceden 2 nodos al log
-Al principio cuando se intenta contactar al data node se hace un random asumiento que las 3 máquinas están disponibles al momento de subir el archivo, por lo que si apagan una máquina y sale justo esa en el random el cliente va a tirar error y terminará su ejecución
