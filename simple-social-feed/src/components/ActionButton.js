import React from "react";
import PropTypes from "prop-types";

const ActionButton = ({ text, onClick, style, onMouseEnter, onMouseLeave }) => {

    const defaultStyle = {
        backgroundColor: "#068932",
        color: "white",
        padding: "12px 12px",
        borderRadius: "5px",
        textAlign: "center",
        textDecoration: "none",
        display: "inline-block",
        ...style
    }

    return (
        <button
            onClick={onClick}
            style={defaultStyle}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
        >
            { text }
        </button>
    )
}

ActionButton.propTypes = {
    text: PropTypes.string,
    onClick: PropTypes.func,
    style: PropTypes.object,
    onMouseEnter: PropTypes.func,
    onMouseLeave: PropTypes.func
}

export default ActionButton;