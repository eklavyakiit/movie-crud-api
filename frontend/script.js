
const API_URL = "http://localhost:8080";

const movieForm = document.getElementById("movie-form");
const movieList = document.getElementById("movie-list");
const message = document.getElementById("message");

const formTitle = document.getElementById("form-title");
const submitButton = document.getElementById("submit-button");
const cancelButton = document.getElementById("cancel-button");
const refreshButton = document.getElementById("refresh-button");

let editingMovieId = null;


// Fetch and display all movies
async function loadMovies() {
    movieList.innerHTML = "<p>Loading movies...</p>";

    try {
        const response = await fetch(`${API_URL}/movies`);

        if (!response.ok) {
            throw new Error("Could not fetch movies");
        }

        const movies = await response.json();

        movieList.innerHTML = "";

        if (!movies || movies.length === 0) {
            movieList.innerHTML =
                '<p class="empty-message">No movies found. Add your first movie!</p>';
            return;
        }

        movies.forEach(function (movie) {
            const card = document.createElement("div");
            card.className = "movie-card";

            const title = document.createElement("h3");
            title.textContent = movie.title;

            const id = document.createElement("span");
            id.className = "movie-id";
            id.textContent = `ID: ${movie.id}`;

            const isbn = document.createElement("p");
            isbn.textContent = `ISBN: ${movie.isbn}`;

            const director = document.createElement("p");
            const directorData = movie.director || {};
            director.textContent =
                `Director: ${directorData.firstname || ""} ${directorData.lastname || ""}`.trim();

            const buttons = document.createElement("div");
            buttons.className = "card-buttons";

            const editButton = document.createElement("button");
            editButton.className = "edit-button";
            editButton.textContent = "Edit";
            editButton.addEventListener("click", function () {
                editMovie(movie);
            });

            const deleteButton = document.createElement("button");
            deleteButton.className = "delete-button";
            deleteButton.textContent = "Delete";
            deleteButton.addEventListener("click", function () {
                deleteMovie(movie.id);
            });

            buttons.appendChild(editButton);
            buttons.appendChild(deleteButton);

            card.appendChild(title);
            card.appendChild(id);
            card.appendChild(isbn);
            card.appendChild(director);
            card.appendChild(buttons);

            movieList.appendChild(card);
        });

    } catch (error) {
        movieList.innerHTML =
            '<p class="empty-message">Could not connect to the Go server. Check whether it is running.</p>';

        console.error(error);
    }
}


// Add a movie or update an existing movie
movieForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const movie = {
        title: document.getElementById("title").value.trim(),
        isbn: document.getElementById("isbn").value.trim(),
        director: {
            firstname: document.getElementById("firstname").value.trim(),
            lastname: document.getElementById("lastname").value.trim()
        }
    };

    const isEditing = editingMovieId !== null;

    const url = isEditing
        ? `${API_URL}/movies/${encodeURIComponent(editingMovieId)}`
        : `${API_URL}/movies`;

    const method = isEditing ? "PUT" : "POST";

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(movie)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || "Request failed");
        }

        message.textContent = isEditing
            ? "Movie updated successfully!"
            : "Movie added successfully!";

        resetForm();
        await loadMovies();

    } catch (error) {
        message.textContent = "Could not save movie. Check the server and console.";
        console.error(error);
    }
});


// Fill the form with the selected movie's information
function editMovie(movie) {
    editingMovieId = movie.id;

    document.getElementById("title").value = movie.title || "";
    document.getElementById("isbn").value = movie.isbn || "";
    document.getElementById("firstname").value =
        movie.director?.firstname || "";
    document.getElementById("lastname").value =
        movie.director?.lastname || "";

    formTitle.textContent = "Edit Movie";
    submitButton.textContent = "Update Movie";
    cancelButton.classList.remove("hidden");

    message.textContent = "";

    document.getElementById("title").focus();
    window.scrollTo({ top: 0, behavior: "smooth" });
}


// Delete a movie
async function deleteMovie(id) {
    const confirmed = confirm("Are you sure you want to delete this movie?");

    if (!confirmed) {
        return;
    }

    try {
        const response = await fetch(
            `${API_URL}/movies/${encodeURIComponent(id)}`,
            { method: "DELETE" }
        );

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || "Could not delete movie");
        }

        message.textContent = "Movie deleted successfully!";
        await loadMovies();

    } catch (error) {
        message.textContent = "Could not delete movie. Check the server and console.";
        console.error(error);
    }
}


// Reset the form back to Add Movie mode
function resetForm() {
    movieForm.reset();
    editingMovieId = null;

    formTitle.textContent = "Add a Movie";
    submitButton.textContent = "Add Movie";
    cancelButton.classList.add("hidden");
}


// Cancel editing
cancelButton.addEventListener("click", function () {
    resetForm();
    message.textContent = "";
});


// Refresh the movie collection
refreshButton.addEventListener("click", loadMovies);


// Load movies when the page opens
loadMovies();
