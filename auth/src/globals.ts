export const typeDefs = `#graphql
   type User {
        id: ID!
        name: String!
        username: String! 
        password: String!
        email: String!
        dateCreated: String!    
   }
   type Query {
        user(email: String!, password: String!)
   }
`
