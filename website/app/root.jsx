import {
  isRouteErrorResponse,
  Outlet,
  Scripts,
} from "react-router";
import "./root.css"

export function Layout({children}) {
    return <html style={{backgroundColor: "#1c1c1c"}}>
                <head>
                    <title>PhotoSync</title>
                    <link rel="icon" type="image/x-icon" href="favicon.ico"></link>
                </head>
                <body>
                    {children}
                    <Scripts/>
                </body>
            </html>
}

export default function App() {
  return <Outlet/>
}

export function HydrateFallback() {
  return <div></div>
}

export function ErrorBoundary({ error }) {
  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
