import { Outlet } from "react-router-dom";
import { CustomSidebar } from "../components/custom/Sidebar";

function RootLayout() {
  return (
    <>
    <CustomSidebar/>
      <Outlet />
    </>
  );
}

export default RootLayout;
