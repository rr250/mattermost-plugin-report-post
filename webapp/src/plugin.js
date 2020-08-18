/* eslint-disable lines-around-comment */
/* eslint-disable no-console */
import {Component} from 'react';
import axios from 'axios';

import {id as pluginId} from './manifest';

// eslint-disable-next-line react/require-optimization
class ReportPlugin extends Component {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction('Report', (postId) => {
            const state = store.getState();
            axios.post('/plugins/' + pluginId + '/getreason', {
                post_id: postId,
                current_user_id: state.user_id,
            }).then().catch((err) => {
                console.log(err);
            });
        });
    }

    uninitialize() {
        console.log(pluginId + '::uninitialize()');
    }
}

export default ReportPlugin;
