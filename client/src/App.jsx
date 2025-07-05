import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { useState } from "react";
import Login from "./pages/Login";
import AdminRoutes from "./routes/AdminRoutes";
import UserRoutes from "./routes/UserRoutes";
import NotFound from "./pages/NotFound";
import Register from "./pages/Register";
import DashboardRedirect from "./pages/DashboardRedirect";
import Notification from "./components/common/Notification";
import "./styles/common/App.css";

const App = () => {
  const [notification, setNotification] = useState({
    message: '',
    type: 'info',
    isVisible: false
  });

  const showNotification = (message, type = 'info') => {
    setNotification({
      message,
      type,
      isVisible: true
    });
  };

  const hideNotification = () => {
    setNotification(prev => ({
      ...prev,
      isVisible: false
    }));
  };

  return (
    <Router>
      <Routes>
        <Route path="/" element={<DashboardRedirect />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        <Route path="/admin/*" element={<AdminRoutes />} />
        <Route path="/user/*" element={<UserRoutes />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
      
      <Notification
        message={notification.message}
        type={notification.type}
        isVisible={notification.isVisible}
        onClose={hideNotification}
      />
    </Router>
  );
};

export default App;