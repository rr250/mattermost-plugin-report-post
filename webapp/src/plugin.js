/* eslint-disable lines-around-comment */
/* eslint-disable no-console */
import {Component} from 'react';
import {Client4} from 'mattermost-redux/client';

import {id as pluginId} from './manifest';

// eslint-disable-next-line react/require-optimization
class ReportPlugin extends Component {
    initialize(registry) {
        registry.registerPostDropdownMenuAction('Report', async (postId) => {
            await fetch(window.location.origin + '/plugins/' + pluginId + '/getreason', Client4.getOptions({
                method: 'post',
                body: JSON.stringify({post_id: postId}),
            }));
        });
    }

    uninitialize() {
        console.log(pluginId + '::uninitialize()');
    }
}

export default ReportPlugin;
