import { index, route, layout} from "@react-router/dev/routes";

export default [
    layout("routes/layout.jsx", [
        index("routes/index.jsx"),
        route("/register", "routes/register.jsx"),
        route("/login", "routes/login.jsx"),
        route("/upload", "routes/upload.jsx"),
        route("/gallery", "routes/gallery.jsx"),
        route("/error", "routes/error.jsx"),
    ])
];
