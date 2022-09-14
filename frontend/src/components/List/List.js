import { useState, useEffect } from 'react';

const List = () => {
    const [elements, setElements] = useState([]);

    const fetchElements = async () => {
        try {
            const response = await fetch(
                'the-new-yorker.json', {
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    }
                }
            );
            const jsonData = await response.json();
            let arrOfLinks = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
            setElements(arrOfLinks);
        }
        catch (err) {
            console.error(err.message);
        }
    }

    useEffect(() => {
        fetchElements();
    }, []);

    return (
        <ul>
            {elements && elements.map((element, idx) => 
                <li key={idx}>
                    {element.target_public_name}: 
                    <a href={element.target_url} target="_blank" rel="noopener noreferrer">{element.target_url.substring(0,70)}...</a>
                </li>
              )
            }
        </ul>
    );
};

export default List;
