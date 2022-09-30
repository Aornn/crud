// db.createUser(
//     {
//         user: "root",
//         pwd: "example",
//         roles: [
//             {
//                 role: "readWrite",
//                 db: "crud123"
//             }
//         ]
//     }
// );
db.createUser(
    {
      user: "root",
      pwd: "password",
      roles: [
        {
          role: "readWrite",
          db: "crud"
        }
      ]
    });
  db.createCollection('users');
//   db.users.insertOne(
//     {
//       name: 'Bill Palmer'
//     }
//   );