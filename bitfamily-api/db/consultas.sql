-- name: BuscarTarea :one
SELECT id, categoria_id, descripcion, recompensa, estado FROM tareas WHERE id = $1;

-- name: GuardarTarea :execrows
UPDATE tareas SET categoria_id = $2, descripcion = $3, recompensa = $4, estado = $5 WHERE id = $1;