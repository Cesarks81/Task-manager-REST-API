# Task Manager API

REST API para gestión de tareas construida con **Go**, **Gin** y **MySQL**.

## Tecnologías

- [Go](https://golang.org/) — lenguaje principal
- [Gin](https://github.com/gin-gonic/gin) — framework HTTP
- [MySQL](https://www.mysql.com/) — base de datos
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) — driver MySQL para Go

## Estructura del proyecto

```
├── main.go               # Punto de entrada, configuración de DB y servidor
├── handlers/
│   └── task.go           # Registro de rutas y lógica HTTP
├── db/
│   └── task_queries.go   # Consultas a la base de datos
├── models/
│   └── tasks.go          # Modelo Task y estados aceptados
├── .env.example          # Ejemplo de variables de entorno
└── go.mod / go.sum       # Dependencias
```

## Requisitos previos

- Go 1.21+
- MySQL corriendo localmente o en remoto
- Una base de datos con la tabla `task`:

```sql
CREATE TABLE task (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(50) NOT NULL,
    createdat   DATETIME NOT NULL,
    completedat DATETIME DEFAULT NULL
);
```

## Configuración

Copia el archivo de ejemplo y rellena tus credenciales:

```bash
cp .env.example .env
```

Contenido de `.env`:

```
DB_DSN=usuario:password@tcp(127.0.0.1:3306)/datago?parseTime=true&loc=Europe%2fMadrid
```

## Instalación y ejecución

```bash
# Clona el repositorio
git clone https://github.com/tu-usuario/task-manager-api.git
cd task-manager-api

# Instala dependencias
go mod download

# Exporta la variable de entorno
export DB_DSN="usuario:password@tcp(127.0.0.1:3306)/datago?parseTime=true&loc=Europe%2fMadrid"

# Arranca el servidor
go run .
```

El servidor escucha por defecto en `http://localhost:8080`.

## Endpoints

| Método   | Ruta                       | Descripción                              |
|----------|----------------------------|------------------------------------------|
| `GET`    | `/api/tasks`               | Obtiene todas las tareas                 |
| `GET`    | `/api/tasks/:id`           | Obtiene una tarea por ID                 |
| `GET`    | `/api/tasks/search`        | Busca tareas por `title` o `status`      |
| `POST`   | `/api/tasks`               | Crea una nueva tarea                     |
| `PUT`    | `/api/tasks/:id`           | Actualiza una tarea existente            |
| `DELETE` | `/api/tasks/:id`           | Elimina una tarea                        |

### Estados válidos

```
new | ongoing | completed
```

### Ejemplos de uso

**Crear una tarea**
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Estudiar Go", "description": "Repasar concurrencia", "status": "new"}'
```

**Obtener todas las tareas**
```bash
curl http://localhost:8080/api/tasks
```

**Buscar por estado**
```bash
curl "http://localhost:8080/api/tasks/search?status=ongoing"
```

**Actualizar una tarea**
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Estudiar Go", "description": "Repasar concurrencia", "status": "completed"}'
```

**Eliminar una tarea**
```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```
