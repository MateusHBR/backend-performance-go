INSERT INTO pessoas (nome, idade, altura, peso) 
VALUES ($1, $2, $3, $4) RETURNING id, nome, idade, altura, peso;

-- INSERT INTO pessoas (nome, idade, altura, peso) 
-- VALUES (:nome, :idade, :altura, :peso) RETURNING *;