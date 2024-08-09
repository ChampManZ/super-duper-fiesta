import React from "react";
import PropTypes from "prop-types";

const ButtonLink = ({ href, text, style, onMouseEnter, onMouseLeave }) => {

    const defaultStyle = {
        backgroundColor: "#f44336",
        color: "white",
        padding: "14px 25px",
        borderRadius: "5px",
        textAlign: "center",
        textDecoration: "none",
        display: "inline-block",
        ...style
    }

    return (
        <a
            href={href}
            style={defaultStyle}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
        >
            { text }
        </a>
    )
}

ButtonLink.propTypes = {
    href: PropTypes.string.isRequired,
    text: PropTypes.node.isRequired,
    style: PropTypes.object,
    onMouseEnter: PropTypes.func,
    onMouseLeave: PropTypes.func
}

export default ButtonLink;
