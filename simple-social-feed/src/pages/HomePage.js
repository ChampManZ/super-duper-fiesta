import React from "react";
import ButtonLink from "../components/ButtonLink";

function HomePage() {

    const registerButton = {
        backgroundColor: "#068932",
        color: "white",
        padding: "12px 12px",
        borderRadius: "5px",
        textAlign: "center",
        textDecoration: "none",
        display: "inline-block",
        marginRight: "10px"
    }

    const loginButton = {
        backgroundColor: "#39424E",
        color: "white",
        padding: "12px 12px",
        borderRadius: "5px",
        textAlign: "center",
        textDecoration: "none",
        display: "inline-block",
        marginRight: "10px"
    }

    return (
        <div className="home-page">
            <ButtonLink href="/register" text="Register" style={registerButton} />
            <ButtonLink href="/login" text="Login" style={loginButton} />
        </div>
    )
}

export default HomePage;