perfecto actualiza specs.
quiero que guardes skills para que el harness trabajé tal cuál lo hicimos. disparando agentes independientes. crea en el harness un hook que siempre obligue siempre que se vaya a enviar un agente independiente (en paralelo o no) a que llame una skill que identifique para el proceso que modelo utilizar dependiendo de complejidad y tamño de la sub tarea.
para specs, una skill que ayude a crearlos facil y rapidamente, de forma muy directa y sin ceremonias como lo recomiendan los archivos del proyecto (lee README y STRUCUTE y TESTING para tenerlos frescos)
deberiamos tener skills para construir backend y frontend por separado y entonces necesito que me organices los archisos iniciales dentro de docs en una estrcutura que permita tener index para saber cada agente que parte necesita (back, front -lo compartido en general de la estrucutra del proyecto para cualquier agente que vaya a tocar codigo)
una skill para hacer busquedas en el codigo con un subagente de menor complejidad... luna?
una skill llamada craft que va a recibir un prompt y con base en eso hacemos toda la spec, y ya como el codigo es tan simple, ese es el plan mismo. esa skill craft puede llamar a las otras skills asi:
spec (debe llevar el happy path y posible errores que puede cometer el usuario para generar los tests)
questions (si quedan pregutnas luego del spec pedir una a una aclaración al programador)
skills de frontend y backend (cada una debe generar sus tests de acuerdo a lo que tenemos definido)
skill run tests funcionales usando el docker
code review (teniendo en cuenta las reglas del repo)
run tests de nuevo