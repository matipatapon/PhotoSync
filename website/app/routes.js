import { index, route, layout} from "@react-router/dev/routes";

export default [
    layout("routes/layout.jsx", [
        index("routes/index.jsx"),
        route("/registration", "routes/registration.jsx"),
        route("/login", "routes/login.jsx"),
        route("/upload", "routes/upload.jsx"),
    ])
];
