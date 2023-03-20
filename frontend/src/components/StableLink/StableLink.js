import { useRef, useState } from 'react';

const StableLink = () => {
    const [inputVisible, setInputVisible] = useState(false);
    const [link, setLink] = useState('');
    const inputRef = useRef(null);

    const [mainButtonHover, setMainButtonHover] = useState(false);
    const [copyButtonHover, setCopyButtonHover] = useState(false);
    const [closeButtonHover, setCloseButtonHover] = useState(false);

    const buttonStyle = {
        borderRadius: '3px',
        backgroundColor: 'white',
        cursor: 'pointer',
        border: '1px solid #ccc',
        padding: '5px 10px',
        marginRight: '5px',
    };

    const closeButtonStyle = {
        ...buttonStyle,
        fontSize: '0.8rem',
    };

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
        inputRef.current.style.color = '#337ab7';
    };

    return (
        <div>
            <button
                onClick={handleClick}
                onMouseEnter={() => setMainButtonHover(true)}
                onMouseLeave={() => setMainButtonHover(false)}
                style={{
                    ...buttonStyle,
                    backgroundColor: mainButtonHover ? '#e6e6e6' : 'white',
                }}
            >
                🔗 Create a stable link to this page
            </button>
            {inputVisible && (
                <div>
                    <input
                        type="text"
                        readOnly
                        value={link}
                        ref={inputRef}
                        style={{
                            borderRadius: '3px',
                            backgroundColor: 'white',
                            marginRight: '5px',
                        }}
                    />
                    <button
                        onClick={copyToClipboard}
                        onMouseEnter={() => setCopyButtonHover(true)}
                        onMouseLeave={() => setCopyButtonHover(false)}
                        style={{
                            ...buttonStyle,
                            backgroundColor: copyButtonHover ? '#e6e6e6' : 'white',
                        }}
                    >
                        Copy
                    </button>
                    <button
                        onClick={handleClose}
                        onMouseEnter={() => setCloseButtonHover(true)}
                        onMouseLeave={() => setCloseButtonHover(false)}
                        style={{
                            ...closeButtonStyle,
                            backgroundColor: closeButtonHover ? '#e6e6e6' : 'white',
                            textDecoration: closeButtonHover ? 'underline' : 'none',
                        }}
                    >
                        X
                    </button>
                </div>
            )}
        </div>
    );
};

export default StableLink;
