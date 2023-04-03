import { useRef, useState } from 'react';

const StableLink = () => {
    const [inputVisible, setInputVisible] = useState(false);
    const [link, setLink] = useState('');
    const inputRef = useRef(null);

    const [mainButtonHover, setMainButtonHover] = useState(false);
    const [copyButtonHover, setCopyButtonHover] = useState(false);
    const [closeButtonHover, setCloseButtonHover] = useState(false);

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

    const handleClick = () => {
        setLink(window.location.href);
        setInputVisible(true);
    };

    const handleClose = () => {
        setInputVisible(false);
    };

    const copyToClipboard = () => {
        navigator.clipboard.writeText(link);
        inputRef.current.focus();
        inputRef.current.style.color = '#3dbbdb';
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
                .stable-link-focus:focus {
                    border: 3px solid black;
                    outline: 3px solid #3DBBDB;
                    padding: 5px;
                    border-radius: 0;
                }
            `}</style>
            <div aria-labelledby="stable-link-label">
                <button
                    onClick={handleClick}
                    onMouseEnter={() => setMainButtonHover(true)}
                    onMouseLeave={() => setMainButtonHover(false)}
                    className="stable-link-focus"
                    id="stable-link-label"
                    style={{
                        ...buttonStyle,
                        backgroundColor: mainButtonHover ? '#6c07ae' : buttonStyle.backgroundColor,
                    }}
                    aria-label="Create a stable link to this page"
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
                            id="stable-link-input-text"
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
                            aria-label="Copy stable link"
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
                            aria-label="Close stable link"
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
