package repository

const createClientQuery = `
	INSERT INTO clients(email) 
	VALUES ($1) 
	ON CONFLICT (email) 
	DO UPDATE SET email = email 
	RETURNING id
`
const createOrderQuery = `
	INSERT INTO orders(client_id) 
	VALUES($1) 
	RETURNING id
`

const getClientQuery = `
	SELECT 
		email
	FROM 
		clients
	WHERE
		id = $1
`
const getOrderQuery = `
	SELECT 
		client_id,
		created_at,
		confirmed_at,
		sended_at,
		delivered_at
	FROM 
		orders
	WHERE
		id = $1
`

const updateOrderQuery = `
	UPDATE 
		orders 
	SET 
		client_id = $2, 
		confirmed_at = $3, 
		sended_at = $4, 
		delivered_at = $5 
	WHERE 
		id = $1
`

const deleteOrderQuery = `
	DELETE FROM orders WHERE id = $1
`
