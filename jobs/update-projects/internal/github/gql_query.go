package github

// GraphQL Query Used
//
// {
//   user(login: "%s") {
//     pinnedItems(first: 6, types: REPOSITORY) {
//       nodes {
//         ... on Repository {
//           name
//           description
//           url
//           languages(first: 20) {
//             edges {
//               node {
//                 name
//               }
//             }
//           }
//           repositoryTopics(first: 20) {
//             edges {
//               node {
//                 topic {
//                   name
//                 }
//               }
//             }
//           }
//           homepageUrl
//         }
//       }
//     }
//   }
// }`

const gqlQueryString = `{"query":"{\n  user(login: \"%s\") {\n    pinnedItems(first: 6, types: REPOSITORY) {\n      nodes {\n        ... on Repository {\n          name\n          description\n          url\n          languages(first: 20) {\n            edges {\n              node {\n                name\n              }\n            }\n          }\n          repositoryTopics(first: 20) {\n            edges {\n              node {\n                topic {\n                  name\n                }\n              }\n            }\n          }\n          homepageUrl\n        }\n      }\n    }\n  }\n}","variables":{}}`
