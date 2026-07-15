import { RouterProvider } from "react-router";
import router from "./routes";
import { ToastContainer } from "react-toastify";
import "react-loading-skeleton/dist/skeleton.css";

function App() {
  return (
    <>
      <RouterProvider router={router} />
      <ToastContainer
        position="top-right"
        autoClose={1500}
        hideProgressBar={false}
      />
    </>
  );
}

export default App;
