document.addEventListener("DOMContentLoaded", function () {
  // registers of users
  document
    .getElementById("registerForm")
    .addEventListener("submit", function (e) {
      e.preventDefault();

      const email = document.getElementById("registerEmail").value;
      const password = document.getElementById("registerPassword").value;

      fetch("http://localhost:8080/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email: email, password: password }),
      })
        .then((response) => response.json())
        .then((data) => {
          alert("User registered successfully!");
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    });

  document.getElementById("loginForm").addEventListener("submit", function (e) {
    e.preventDefault(); // prevent default comportament of forms (refresh page)

    const email = document.getElementById("loginEmail").value;
    const password = document.getElementById("loginPassword").value;

    // Requisition POST to the API login
    fetch("http://localhost:8080/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email, password: password }),
    })
      .then(async (response) => {
        if (!response.ok) {
          const text = await response.text();
            throw new Error(text);
        }
        return response.json()
      })

      .then(data => {
        console.log(Object.values(data))
        console.log(Object.keys(data))
        if (data.success) {
          // Save the user id in localStorage
          localStorage.setItem("userID", data.userID);
          // Redirect to the home page if login is succeful
          alert("Login sucess");
          window.location.href = "/front/decks.html";
        } else {
          alert("Login failed!");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("Falha no login, tente novamente.");
      });
  });

//   document.getElementById("deckForm").addEventListener("submit", function (e) {
//     e.preventDefault();

//     const deckName = document.getElementById("deckName").value;

//     fetch("http://localhost:8080/decks", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({ name: deckName }),
//     })
//       .then((response) => response.json())
//       .then((data) => {
//         alert("Deck created successfully!");
//       })
//       .catch((error) => {
//         console.error("Error:", error);
//       });
//   });

//   document
//     .getElementById("winrateForm")
//     .addEventListener("submit", function (e) {
//       e.preventDefault();

//       const userDeck = document.getElementById("userDeck").value;
//       const opponentDeck = document.getElementById("opponentDeck").value;
//       const wins = document.getElementById("wins").value;
//       const losses = document.getElementById("losses").value;

//       fetch("http://localhost:8080/matches", {
//         method: "POST",
//         headers: {
//           "Content-Type": "application/json",
//         },

//         body: JSON.stringify({ userDeck, opponentDeck, wins, losses }),
//       })
//         .then((response) => response.json())
//         .then((data) => {
//           alert("Win/Loss rate updated!");
//         })
//         .catch((error) => {
//           console.error("Error:", error);
//         });
//     });
});
