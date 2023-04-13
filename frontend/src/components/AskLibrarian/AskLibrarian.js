const AskLibrarian = () => {
    return (
        <div className="ask-librarian">
            <h2>Need Help?</h2>
            <h3>
                {ASK_LIBRARIAN_TEXT}
            </h3>
            <div>
                <p>
                    Use <a href="https://library.nyu.edu/ask/" target="_blank" rel="noreferrer">{ASK_LIBRARIAN_TEXT}</a> or the &quot;Chat with Us&quot; icon at the bottom right corner for any question you have about the Libraries&apos; services.
                </p>
                <p>
                    Visit our <a href="https://guides.nyu.edu/online-tutorials/finding-sources" target="_blank" rel="noreferrer">online tutorials</a> for tips on searching the catalog and getting library resources.
                </p>
                <h3>Additional Resources</h3>
                <ul>
                    <li>Use <a href="https://ezborrow.reshare.indexdata.com/" target="_blank" rel="noreferrer">EZBorrow</a> or <a href="https://library.nyu.edu/services/borrowing/from-non-nyu-libraries/interlibrary-loan/" target="_blank" rel="noreferrer">InterLibrary Loan (ILL)</a> for materials unavailable at NYU</li>
                    <li>Discover subject specific resources using <a href="http://guides.nyu.edu" target="_blank" rel="noreferrer">expert curated research guides</a></li>
                    <li>Explore the <a href="https://library.nyu.edu/services/" target="_blank" rel="noreferrer">complete list of library services</a></li>
                    <li>Search <a href="https://www.worldcat.org/search?qt=worldcat_org_all" target="_blank" rel="noreferrer">WorldCat</a> for items in nearby libraries</li>
                    <li><a href="https://nyu.qualtrics.com/jfe/form/SV_a30SvEFe8f1yXQh" target="_blank" rel="noreferrer">Problem with an e-resource? Fill out this form</a></li>
                    <li><a href="https://nyu.qualtrics.com/jfe/form/SV_blQ3OFOew9vl6Pb?Source=NYU" target="_blank" rel="noreferrer">Submit feedback about GetIt</a></li>
                </ul>
            </div>
        </div>
    );
};

export const ASK_LIBRARIAN_TEXT = 'Ask a Librarian';
export const ASK_LIBRARIAN_URL = 'https://library.nyu.edu/ask/';

export default AskLibrarian;
