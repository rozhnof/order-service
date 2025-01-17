package repository

const createOrderQuery = `
INSERT INTO orders(id) VALUES ($1)
`
