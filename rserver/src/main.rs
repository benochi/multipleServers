use actix_web::{error::ErrorInternalServerError, web, App, HttpResponse, HttpServer};
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

#[derive(Serialize, Deserialize, Clone)]
struct Todo {
    id: u32,
    task: String,
    completed: bool,
}
#[derive(Serialize, Deserialize)]
struct NewTodo {
    task: String,
    completed: bool,
}

#[derive(Serialize, Deserialize)]
struct Todos {
    todos: Vec<Todo>,
}

fn read_todos() -> Result<Todos, actix_web::Error> {
    let mut path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    path.push("..");
    path.push("db.json");

    let contents = fs::read_to_string(&path).map_err(ErrorInternalServerError)?;

    serde_json::from_str(&contents).map_err(ErrorInternalServerError)
}

fn write_todos(todos: &Todos) -> Result<(), actix_web::Error> {
    let mut path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    path.push("..");
    path.push("db.json");

    let todos_json = serde_json::to_string(&todos).map_err(ErrorInternalServerError)?;

    fs::write(path, todos_json).map_err(ErrorInternalServerError)
}

async fn get_todos() -> Result<HttpResponse, actix_web::Error> {
    read_todos().map(|todos| HttpResponse::Ok().json(todos))
}

async fn put_todo(
    id: web::Path<u32>,
    input_todo: web::Json<Todo>,
) -> Result<HttpResponse, actix_web::Error> {
    let mut todos = read_todos()?;
    let mut updated = false;

    for todo in &mut todos.todos {
        if todo.id == *id {
            *todo = input_todo.clone();
            updated = true;
            break;
        }
    }

    if updated {
        write_todos(&todos)?;
        Ok(HttpResponse::Ok().json(input_todo.into_inner()))
    } else {
        Ok(HttpResponse::NotFound().finish())
    }
}

async fn post_todo(input_todo: web::Json<NewTodo>) -> Result<HttpResponse, actix_web::Error> {
    let mut todos = read_todos()?;
    let new_id = todos.todos.iter().map(|todo| todo.id).max().unwrap_or(0) + 1;

    let new_todo = Todo {
        id: new_id,
        task: input_todo.task.clone(),
        completed: input_todo.completed,
    };

    todos.todos.push(new_todo.clone());
    write_todos(&todos)?;

    Ok(HttpResponse::Created().json(new_todo))
}
async fn delete_todo(id: web::Path<u32>) -> Result<HttpResponse, actix_web::Error> {
    let mut todos = read_todos()?;
    let initial_len = todos.todos.len();

    todos.todos.retain(|todo| todo.id != *id);

    if todos.todos.len() < initial_len {
        write_todos(&todos)?;
        Ok(HttpResponse::Ok().finish())
    } else {
        Ok(HttpResponse::NotFound().finish())
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/todos", web::get().to(get_todos))
            .route("/todos/{id}", web::put().to(put_todo))
            .route("/todos", web::post().to(post_todo))
            .route("/todos/{id}", web::delete().to(delete_todo))
    })
    .bind("127.0.0.1:6003")?
    .run()
    .await
}
