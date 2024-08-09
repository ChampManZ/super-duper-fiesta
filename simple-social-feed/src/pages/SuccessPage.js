import React from "react";
import ButtonLink from "../components/ButtonLink";

function SuccessPage() {

    return (
        <div className="success-page">
            <h2>Thank you for registering!</h2>
            <p>Your account has been created successfully.</p>
            <ButtonLink 
                href="/" 
                text="Go back to the registration page" 
                onMouseEnter={e => e.currentTarget.style.backgroundColor = 'red'} 
                onMouseLeave={e => e.currentTarget.style.backgroundColor = '#f44336'} 
            />
        </div>
    )
}

export default SuccessPage;
