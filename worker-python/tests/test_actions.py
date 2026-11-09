from robotdatahub_worker.validators.actions import validate_actions

def test_valid():
    assert validate_actions([{'frame_index': 0, 'action': [0.] * 14}]) == []

def test_nan():
    issues = validate_actions([{'frame_index': 10, 'action': [float('nan')] + [0.] * 13}])
    assert issues[0]['code'] == 'NON_FINITE'

def test_dimension():
    assert validate_actions([{'frame_index': 0, 'action': [1]}])[0]['code'] == 'ACTION_DIM'
