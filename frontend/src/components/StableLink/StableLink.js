import { useRef, useState } from 'react';

const StableLink = () => {
    const [inputVisible, setInputVisible] = useState(false);
    const [link, setLink] = useState('');
    const inputRef = useRef(null);
    const mainButtonRef = useRef(null);
    const closeButtonRef = useRef(null);

    const [mainButtonHover, setMainButtonHover] = useState(false);
    const [copyButtonHover, setCopyButtonHover] = useState(false);
    const [closeButtonHover, setCloseButtonHover] = useState(false);
    const [mainButtonFocus, setMainButtonFocus] = useState(false);

    const buttonStyle = {
        backgroundColor: '#57068c',
        color: '#fff',
        cursor: 'pointer',
        border: '1px solid #ccc',
        padding: '5px 10px',
        marginRight: '5px',
        marginBottom: '5px',
    };

    const closeButtonStyle = {
        ...buttonStyle,
        fontSize: '0.8rem',
        color: '#1C2127',
    };

    const copyButtonStyle = {
        ...buttonStyle,
        padding: '2px 10px',
    }

    const focusedStyle = {
        border: '3px solid black',
        outline: '3px solid #3DBBDB',
        padding: '5px',
        borderRadius: '0',
    };

    const handleClick = () => {
        setLink(window.location.href);
        setInputVisible(true);
    };

    const handleClose = () => {
        setInputVisible(false);
        mainButtonRef.current.focus();
    };

    const copyToClipboard = () => {
        navigator.clipboard.writeText(link);
        inputRef.current.focus();
        inputRef.current.style.color = '#57068c';
    };

    return (
        <>
            <style>{`
        .link-icon::before {
          content: "\\e157";
          font-family: "Material Symbols Sharp";
          color: #fff;
          margin-right: 5px;
          vertical-align: middle;
        }
      `}</style>
            <div aria-labelledby="stable-link-label">
                <button
                    onClick={handleClick}
                    onMouseEnter={() => setMainButtonHover(true)}
                    onMouseLeave={() => setMainButtonHover(false)}
                    onFocus={() => setMainButtonFocus(true)}
                    onBlur={() => setMainButtonFocus(false)}
                    id="stable-link-label"
                    ref={mainButtonRef}
                    style={{
                        ...buttonStyle,
                        ...mainButtonFocus ? focusedStyle : {},
                        backgroundColor: mainButtonHover ? '#6c07ae' : buttonStyle.backgroundColor,
                    }}
                    aria-label="stable-link-label"
                >
                    <span className="link-icon"></span>Create a stable link to this page
                </button>
                {inputVisible && (
                    <div>
                        <input
                            type="text"
                            readOnly
                            value={link}
                            ref={inputRef}
                            id="stable-link-text"
                            style={{
                                borderRadius: '3px',
                                backgroundColor: 'white',
                                marginRight: '5px',
                            }}
                            aria-labelledby="stable-link-text"
                        />
                        <button
                            onClick={copyToClipboard}
                            onMouseEnter={() => setCopyButtonHover(true)}
                            onMouseLeave={() => setCopyButtonHover(false)}
                            id="copy-stable-link"
                            style={{
                                ...copyButtonStyle,
                                color: '#1C2127',
                                backgroundColor: copyButtonHover ? '#e6e6e6' : 'white',
                            }}
                            aria-label="copy-stable-link"
                        >
                            Copy
                        </button>
                        <button
                            onClick={handleClose}
                            onMouseEnter={() => setCloseButtonHover(true)}
                            onMouseLeave={() => setCloseButtonHover(false)}
                            id="close-stable-link"
                            style={{
                                ...closeButtonStyle,
                                backgroundColor: closeButtonHover ? '#e6e6e6' : 'white',
                                textDecoration: closeButtonHover ? 'underline' : 'none',
                            }}
                            aria-label="close-stable-link"
                            ref={closeButtonRef}
                        >
                            X
                        </button>
                    </div>
                )}
            </div>
        </>
    );
};

export default StableLink;
