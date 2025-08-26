import { index, route, layout} from "@react-router/dev/routes";

export default [
    layout("routes/layout.jsx", [
        index("routes/upload.jsx")
        // index("routes/authentication/index.jsx"),
        // route("/registration", "routes/authentication/registration.jsx"),
        // route("/login", "routes/authentication/login.jsx"),
    ])
];
