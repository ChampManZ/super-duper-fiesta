import './App.css';
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import RegisterForm from './components/UserForm';
import SuccessPage from './components/SuccessPage';

function App() {
  return (
    <Router>
      <div className="App">
        <h1>UrMessage - Place to leave your message</h1>
        <Routes>
          <Route exact path="/" element={<RegisterForm />} />
          <Route path="/success" element={<SuccessPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;

// References: https://reactrouter.com/en/main/upgrading/v5#upgrade-all-switch-elements-to-routes
// React Router v.6 has removed the <Switch> component. Instead, use <Routes> and <Route> components.