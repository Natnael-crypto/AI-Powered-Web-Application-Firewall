import { createBrowserRouter, createRoutesFromElements, Route, RouterProvider } from "react-router-dom";
import RootLayout from "./layout/RootLayout";
import { SidebarProvider } from "./components/ui/sidebar";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<RootLayout />}>
    </Route>
  )
);

function App() {
  return (
    <>
    <SidebarProvider>
      <RouterProvider router={router} />
      </SidebarProvider>
    </>

  );
}

export default App;

