import React from "react";
import PropTypes from "prop-types";

function MessageFeed({ username, postTime, message }) {

    const defaultStyle = {
        border: "2px solid black",
        borderRadius: "10px",
        padding: "20px",
        maxWidth: "600px",
        margin: "20px auto",
        backgroundColor: "#f9f9f9",
        textAlign: "left"
    }

    const usernameStyle = {
        fontWeight: "bold",
        fontSize: "18px",
    }

    const postTimeStyle = {
        fontSize: "14px",
        color: "gray",
        marginBottom: "10px",
    }

    const messageStyle = {
        fontSize: "16px",
    }

    return (
        <div style={defaultStyle}>
            <div style={usernameStyle}>{username}</div>
            <div style={postTimeStyle}>{new Date(postTime).toLocaleString()}</div>
            <div style={messageStyle}>{message}</div>
        </div>
    )
}

MessageFeed.propTypes = {
    username: PropTypes.string,
    postTime: PropTypes.string,
    message: PropTypes.string,
}

export default MessageFeed;