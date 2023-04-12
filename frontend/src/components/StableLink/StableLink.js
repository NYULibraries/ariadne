import { useRef, useState } from 'react';

const StableLink = () => {
    // State for handling the copy status: null, success, or error
    const [copyStatus, setCopyStatus] = useState(null);
    const mainButtonRef = useRef(null);

    const [mainButtonHover, setMainButtonHover] = useState(false);
    const [mainButtonFocus, setMainButtonFocus] = useState(false);

    // Function to handle the click event on the main button
    const handleClick = async () => {
        const link = window.location.href;

        try {
            // Attempt to use the Clipboard API to copy the link
            await navigator.clipboard.writeText(link);
            setCopyStatus('success');
        } catch (err) {
            console.error('Error copying text using Clipboard API: ', err);
            setCopyStatus('error');
        }

        // Reset the copy status to null after 5 seconds
        if (copyStatus !== 'error') {
            setTimeout(() => {
                setCopyStatus(null);
            }, 5000);
        }
    };

    const errorMessage = 'Error copying link. Please try again or manually copy the link.';


    const buttonStyle = {
        backgroundColor: '#57068c',
        color: '#fff',
        cursor: 'pointer',
        border: '1px solid #ccc',
        padding: '5px 10px',
        marginRight: '5px',
        marginBottom: '5px',
    };

    const focusedStyle = {
        border: '3px solid black',
        outline: '3px solid #3DBBDB',
        padding: '5px',
        borderRadius: '0',
    };

    const copiedTextStyle = {
        display: 'inline-block',
        lineHeight: '30px', // adjust this value according to your button height
        padding: '5px 10px',
        marginRight: '5px',
        marginBottom: '5px',
    };

    return (
        <div aria-labelledby="stable-link-label">
            {copyStatus !== 'success' ? (
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
                    {copyStatus === 'error'
                        ? errorMessage
                        : 'Copy a stable link to this page'}
                </button>
            ) : (
                <span style={{ ...copiedTextStyle }}>Copied!</span>
            )}
            {/* Aria-live region for status messages */}
            {/* About aria-live: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/ARIA_Live_Regions */}
            <div
                aria-live="polite"
                style={{
                    position: 'absolute',
                    width: 0,
                    height: 0,
                    overflow: 'hidden',
                    whiteSpace: 'nowrap',
                }}
            >
                {copyStatus === 'success'
                    ? 'Copied! The stable link has been copied to your clipboard.'
                    : copyStatus === 'error'
                        ? errorMessage
                        : ''}
            </div>
        </div>
    );

};


export default StableLink;