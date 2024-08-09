import './App.css';
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import RegisterPage from './pages/RegisterPage';
import SuccessPage from './pages/SuccessPage';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import FeedPage from './pages/FeedPage';

function App() {
  return (
    <Router>
      <div className="App">
        <h1>UrMessage - Place to leave your message</h1>
        <Routes>
          <Route exact path="/" element={<HomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/success" element={<SuccessPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/feed" element={<FeedPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;

// References: https://reactrouter.com/en/main/upgrading/v5#upgrade-all-switch-elements-to-routes
// React Router v.6 has removed the <Switch> component. Instead, use <Routes> and <Route> components.