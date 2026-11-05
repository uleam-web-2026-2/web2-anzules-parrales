CREATE TABLE categorias (
    id bigserial PRIMARY KEY,
    nombre text NOT NULL
);

CREATE TABLE tareas (
    id bigserial PRIMARY KEY,
    categoria_id bigint NOT NULL,
    descripcion text NOT NULL,
    recompensa bigint NOT NULL,
    estado text NOT NULL
);