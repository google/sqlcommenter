suite('timing-utilities', function() {
  test('normalize timing input', function() {
    assert.equal(normalizeTimingInput(1).duration, 1);
    assert.equal(normalizeTimingInput(1)._easingFunction(0.2), 0.2);
    assert.equal(normalizeTimingInput(undefined).duration, 0);
  });
  test('calculating active duration', function() {
    assert.equal(calculateActiveDuration({duration: 1000, playbackRate: 4, iterations: 20}), 5000);
    assert.equal(calculateActiveDuration({duration: 500, playbackRate: 0.1, iterations: 300}), 1500000);
  });
  test('conversion of timing functions', function() {
    function assertTimingFunctionsEqual(tf1, tf2, message) {
      for (var i = 0; i <= 1; i += 0.1) {
        assert.closeTo(tf1(i), tf2(i), 0.01, message);
      }
    }

    assertTimingFunctionsEqual(
        toTimingFunction('ease-in-out'),
        toTimingFunction('eAse\\2d iN-ouT'),
        'Should accept arbitrary casing and escape chararcters');

    var f = toTimingFunction('ease');
    var g = toTimingFunction('cubic-bezier(.25, 0.1, 0.25, 1.0)');
    assertTimingFunctionsEqual(f, g, 'ease should map onto preset cubic-bezier');
    assert.closeTo(f(0.1844), 0.2601, 0.01);
    assert.closeTo(g(0.1844), 0.2601, 0.01);
    assert.equal(f(0), 0);
    assert.equal(f(1), 1);
    assert.equal(g(0), 0);
    assert.equal(g(1), 1);

    f = toTimingFunction('cubic-bezier(0, 1, 1, 0)');
    assert.closeTo(f(0.104), 0.392, 0.01);

    function assertIsLinear(easing) {
      var f = toTimingFunction(easing);
      for (var i = 0; i <= 1; i += 0.1) {
        assert.equal(f(i), i, easing + ' is linear');
      }
    }

    assertIsLinear('cubic-bezier(.25, 0.1, 0.25, 1.)');
    assertIsLinear('cubic-bezier(0, 1, -1, 1)');
    assertIsLinear('an elephant');
    assertIsLinear('cubic-bezier(-1, 1, 1, 1)');
    assertIsLinear('cubic-bezier(1, 1, 1)');

    f = toTimingFunction('steps(10, end)');
    assert.equal(f(0), 0);
    assert.equal(f(0.09), 0);
    assert.equal(f(0.1), 0.1);
    assert.equal(f(0.25), 0.2);
  });
  test('calculating phase', function() {
    // calculatePhase(activeDuration, localTime, timing);
    assert.equal(calculatePhase(1000, 100, {delay: 0}), PhaseActive);
    assert.equal(calculatePhase(1000, 100, {delay: 200}), PhaseBefore);
    assert.equal(calculatePhase(1000, 2000, {delay: 200}), PhaseAfter);
    assert.equal(calculatePhase(1000, null, {delay: 200}), PhaseNone);
  });
  test('calculating active time', function() {
    // calculateActiveTime(activeDuration, fillMode, localTime, phase, delay);
    assert.equal(calculateActiveTime(1000, 'forwards', 100, PhaseActive, 0), 100);
    assert.equal(calculateActiveTime(1000, 'forwards', 100, PhaseBefore, 200), null);
    assert.equal(calculateActiveTime(1000, 'both', 100, PhaseBefore, 200), 0);
    assert.equal(calculateActiveTime(1000, 'forwards', 500, PhaseActive, 200), 300);
    assert.equal(calculateActiveTime(1000, 'forwards', 1100, PhaseAfter, 200), 1000);
    assert.equal(calculateActiveTime(1000, 'none', 1100, PhaseAfter, 200), null);
    assert.equal(calculateActiveTime(Infinity, 'both', 5000000, PhaseActive, 2000000), 3000000);
    assert.equal(calculateActiveTime(Infinity, 'both', 50000, PhaseBefore, 2000000), 0);
  });
  test('calculating scaled active time', function() {
    // calculateScaledActiveTime(activeDuration, activeTime, startOffset, timingInput);
    assert.equal(calculateScaledActiveTime(1000, 200, 300, {playbackRate: 1.5}), 600);
    assert.equal(calculateScaledActiveTime(1000, 200, 300, {playbackRate: -4}), 3500);
    assert.equal(calculateScaledActiveTime(Infinity, 400, 200, {playbackRate: 1}), 600);
    assert.equal(calculateScaledActiveTime(Infinity, 400, 200, {playbackRate: -4}), Infinity);
  });
  test('calculating iteration time', function() {
    // calculateIterationTime(iterationDuration, repeatedDuration, scaledActiveTime, startOffset, timingInput);
    assert.equal(calculateIterationTime(500, 5000, 600, 100, {iterations: 10, iterationStart: 0}), 100);
    assert.equal(calculateIterationTime(500, 5000, Infinity, 100, {iterations: 10, iterationStart: 0}), 500);
    assert.equal(calculateIterationTime(500, 5000, 5100, 100, {iterations: 3.2, iterationStart: 0.8}), 500);
  });
  test('calculating current iteration', function() {
    // calculateCurrentIteration(iterationDuration, iterationTime, scaledActiveTime, timingInput);
    assert.equal(calculateCurrentIteration(1000, 400, 4400, {iterations: 50, iterationStart: 0.8}), 4);
    assert.equal(calculateCurrentIteration(1000, 1000, 4400, {iterations: 50.2, iterationStart: 0.8}), 50);
  });
  test('calculating transformed time', function() {
    // calculateTransformedTime(currentIteration, iterationDuration, iterationTime, timingInput);
    assert.equal(calculateTransformedTime(4, 1000, 200, {_easingFunction: function(x) { return x; }, direction: 'normal'}), 200);
    assert.equal(calculateTransformedTime(4, 1000, 200, {_easingFunction: function(x) { return x; }, direction: 'reverse'}), 800);
    assert.closeTo(calculateTransformedTime(4, 1000, 200, {_easingFunction: function(x) { return x * x; }, direction: 'reverse'}), 640, 0.0001);
    assert.closeTo(calculateTransformedTime(4, 1000, 600, {_easingFunction: function(x) { return x * x; }, direction: 'alternate'}), 360, 0.0001);
    assert.closeTo(calculateTransformedTime(3, 1000, 600, {_easingFunction: function(x) { return x * x; }, direction: 'alternate'}), 160, 0.0001);
    assert.closeTo(calculateTransformedTime(4, 1000, 600, {_easingFunction: function(x) { return x * x; }, direction: 'alternate-reverse'}), 160, 0.0001);
    assert.closeTo(calculateTransformedTime(3, 1000, 600, {_easingFunction: function(x) { return x * x; }, direction: 'alternate-reverse'}), 360, 0.0001);
  });
  test('EffectTime', function() {
    var timing = normalizeTimingInput({duration: 1000, iterations: 4, iterationStart: 0.5, easing: 'linear', direction: 'alternate', delay: 100, fill: 'forwards'});
    var timing2 = normalizeTimingInput({duration: 1000, iterations: 4, iterationStart: 0.5, easing: 'ease', direction: 'alternate', delay: 100, fill: 'forwards'});
    var effectTF = effectTime(timing);
    var effectTF2 = effectTime(timing2);
    assert.equal(effectTF(0), null);
    assert.equal(effectTF(100), 0.5);
    assert.closeTo(effectTF2(100), 0.8, 0.005);
    assert.equal(effectTF(600), 1);
    assert.closeTo(effectTF2(600), 1, 0.005);
    assert.equal(effectTF(700), 0.9);
    assert.closeTo(effectTF2(700), 0.99, 0.005);
    assert.equal(effectTF(1600), 0);
    assert.closeTo(effectTF2(1600), 0, 0.005);
    assert.equal(effectTF(4000), 0.4);
    assert.closeTo(effectTF2(4000), 0.68, 0.005);
    assert.equal(effectTF(4100), 0.5);
    assert.closeTo(effectTF2(4100), 0.8, 0.005);
    assert.equal(effectTF(6000), 0.5);
    assert.closeTo(effectTF2(6000), 0.8, 0.005);
  });
});
